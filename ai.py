# Import libraries
from langchain.vectorstores.cassandra import Cassandra
from langchain.indexes.vectorstore import VectorStoreIndexWrapper
from langchain.llms import OpenAI
from langchain.embeddings import OpenAIEmbeddings
from langchain.schema import Document

from cassandra.cluster import Cluster
from cassandra.auth import PlainTextAuthProvider

from datasets import load_dataset, Dataset

import tiktoken

import sqlite3
import pandas as pd
import threading

import time

# Set up env variables
import os

from dotenv import load_dotenv


load_dotenv()  # take environment variables from .env.


ASTRA_DB_BUNDLE_PATH =os.getenv("ASTRA_DB_BUNDLE_PATH")
ASTRA_DB_TOKEN = os.getenv("TOKEN")
ASTRA_DB_CLIENT_ID = os.getenv("CLIENT_ID")
ASTRA_DB_CLIENT_SECRET = os.getenv("SECRET")
ASTRA_DB_KEYSPACE = os.getenv("KEYSPACE")
OPENAI_KEY= os.getenv("OPENAI_KEY")


#define function for getting tokens from a string
def num_tokens_from_string(string: str, encoding_name: str) -> int:
    """Returns the number of tokens in a text string."""
    encoding = tiktoken.get_encoding(encoding_name)
    num_tokens = len(encoding.encode(string))
    return num_tokens

# Import listing data
def get_listings():
    """Returns a list of all listings."""
    # Connect to the database.
    db = sqlite3.connect("db/listings.db")
    cursor = db.cursor()

    # Get all listings from the database.
    cursor.execute("SELECT * FROM listings")

    # Create a list of all listings.
    listings = []
    for row in cursor.fetchall():
        listings.append(row)

    # Close the connection to the database.
    db.close()

    # Return the list of listings.
    return listings

# Config Astra
cloud_config = {
    "secure_connect_bundle":  ASTRA_DB_BUNDLE_PATH
}
auth_provider = PlainTextAuthProvider(ASTRA_DB_CLIENT_ID, ASTRA_DB_CLIENT_SECRET)
cluster = Cluster(cloud=cloud_config, auth_provider=auth_provider)
astra_session = cluster.connect()

llm = OpenAI(openai_api_key=OPENAI_KEY)
myEmbedding = OpenAIEmbeddings(openai_api_key=OPENAI_KEY)

print("Configuration complete")

# Create Cassandra Store and table if it doesn't exist
listingCassandraStore = Cassandra(
    embedding=myEmbedding,
    session=astra_session,
    keyspace=ASTRA_DB_KEYSPACE,
    table_name="listings"
)

print("Cassandra Store created")

# retrieve listings from sqlite
items = get_listings()

print("Listings fetched")

# create dataframe and clean it up
df = pd.DataFrame(items, columns=["id", "title","price", "link"])
df["price"] = df["price"].apply(lambda p: int(p.replace('$', '').replace(',', '')))
df = df[["title", "price", "link"]]
df = df[df["price"] != 0]

# sort decending price
df = df.sort_values("price", ascending=False)

print("Listings cleaned")

# Set Size of chunks to publish. Larger chunks make embedding much faster 
chunkSize = 100

# OpenAI will limit tokens per minute, chunking will not help you here and you may have to do at separate times
# todo: replace notebook with a script that counts token use and cools down
sizeToEmbed = len(df)
startIndex = 0
currentTokenCount = 0
tokenLimit = 1000000

textListings = []
for i, row in df.iterrows():
    textListings.append(row["title"] + ": $"+str(row["price"]) +"\n")

print("Listings converted to text successfully")

print("Beginning embedding with chunk  of size: ",str(chunkSize), " At: ", str(startIndex), " Total Size: ", str(sizeToEmbed),"/",str(len(df)))
for i in range(startIndex, sizeToEmbed, chunkSize):
    i_end = (i+chunkSize) % (sizeToEmbed-1)
    chunk = textListings[i:i_end]
    chunk_flat = ''.join(chunk)
    currentTokenCount += num_tokens_from_string(chunk_flat,"cl100k_base")

    # if the tokens overflow wait 60 seconds after all pushing to db is complete
    if currentTokenCount > tokenLimit:
        print("Hit token limit: ", currentTokenCount, " tokens")
        print(" Total Size: ", str(sizeToEmbed),"/",str(len(df)))
        # task.join()
        time.sleep(60)

        
        print("Waiting complete")
        
        currentTokenCount = 0

    # thread the publishing of each chunk
    print("pushing chunk...")
    val = listingCassandraStore.add_documents(documents=[Document(page_content=chunk_flat)])
    print("Added: ", val, " documents")
    # task = threading.Thread(target=, args=([Document(page_content=chunk_flat)]))
    # task.start()

print(f"\n Embedded {sizeToEmbed}/{len(df)} listings")

# Preform a query
# vectorIndex = VectorStoreIndexWrapper(vectorstore=listingCassandraStore)

# query = "What is a normal price for a Mercedes e350?"
# answer = vectorIndex.query(question=query, llm=llm).strip()

# print(answer)

# print("Docs by relevance")
# for doc, score in listingCassandraStore.similarity_search_with_score(query, k=4):
#     print("Score:\t",score,"\n",doc)