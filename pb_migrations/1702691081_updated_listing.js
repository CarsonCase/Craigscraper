/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g82ttzu2p9w3y9v")

  collection.indexes = [
    "CREATE UNIQUE INDEX `idx_N4Fc43J` ON `listing` (\n  `title`,\n  `price`,\n  `link`,\n  `created`,\n  `updated`\n)"
  ]

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("g82ttzu2p9w3y9v")

  collection.indexes = []

  return dao.saveCollection(collection)
})
