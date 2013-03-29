# dev.engine.likeit.nexum.ws

Internal tools for analyzing LikeIt data.

## Authors

* José Carlos Nieto <<jose.carlos@menteslibres.org>>

## Bots

### ig-pull-worker

This worker pulls data from both the MySQL view `ig_media_likes` and the
Instagram API, for each of the pulled items a new entry is added to the
`photos` collection of a MongoDB database.

### ig-rank-worker

This worker uses data from the `photos` collection of a MongoDB database
and generates a ranking weight (photo rank). At this time this is a dummy rank
that just applies a factor to specific photo properties.

## API

This HTTP-JSON interface allows the user to request special data gathered by
LikeIt, such as popular photos.

### Prefix

While in development you can use the following API prefix:

`http://dev.engine.likeit.nexum.ws:9192/`

### Response format

This is the general format of an API response:

```
[{}, null]
```

A response is an array of two elements, the first element contains an object or
`null`, the second response contains a error message string or `null`.

An example response with no errors:

```
[{
  "data": [
    {
      "foo": 1
    },
    {
      "foo": 2
    }
  ]
}, null]
```

An example response with errors:

```
[
  null,
  "Something wrong happened!"
]
```

There is no guarantee these two elements are mutually exclusive, it's up to the
user to decide when to use one or the other.

### Endpoints

#### /api/v1/featured/list

Returns a list of the top 20 featured photos.

```
# Top 20 photos.
curl "dev.engine.likeit.nexum.ws:9192/api/v1/featured/list"
```

```
# Top 3 photos.
curl "dev.engine.likeit.nexum.ws:9192/api/v1/featured/list" -d "limit=3"
```
