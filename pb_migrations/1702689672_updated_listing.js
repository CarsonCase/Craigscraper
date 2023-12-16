/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("6engww6nlnv6d58")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "xg4sytet",
    "name": "city",
    "type": "url",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "exceptDomains": null,
      "onlyDomains": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("6engww6nlnv6d58")

  // remove
  collection.schema.removeField("xg4sytet")

  return dao.saveCollection(collection)
})
