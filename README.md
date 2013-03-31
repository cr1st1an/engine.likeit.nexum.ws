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

There is no guarantee that these two elements are mutually exclusive, it's up to
the user to decide when to use one or the other.

### Endpoints

#### /api/v1/photo/set_handpicked_rank

Sets the handpicked value of the photo ranking. The photo ID must exists in the photo listing.

```
# Setting rank for photo 410411077197573167_312334042
curl "dev.engine.likeit.nexum.ws:9192/api/v1/photo/set_handpicked_rank" -d "handpicked_rank=123&id=410411077197573167_312334042"
```

While this operation takes place instantaneously, the general rankings (those
of the /api/v1/featured endpoint) may take up to five minutes to be updated.

#### /api/v1/photos/list

Returns a variable list of photos.

```
# First page. 42 photos per page.
curl "dev.engine.likeit.nexum.ws:9192/api/v1/photos/list" -d "limit=42"
```

```
# Second page. 42 photos per page.
curl "dev.engine.likeit.nexum.ws:9192/api/v1/photos/list" -d "limit=42&page=2"
```

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

```
[
   {
      "data":[
         {
            "info":{
               "_id":"51558cfbb79e40b00067f780",
               "created_time":1352323376,
               "hand_picked":0,
               "id":"319634107029710721_10206720",
               "ig_comments_count":28963,
               "ig_likes_count":258193,
               "likeit_count":2,
               "modified":1364574667,
               "rank":304067
            },
            "photo":{
               "_id":"51558908b79e40b00067f4f8",
               "attribution":null,
               "caption":{
                  "created_time":"1352323383",
                  "from":{
                     "full_name":"Barack Obama",
                     "id":"10206720",
                     "profile_picture":"http://images.instagram.com/profiles/profile_10206720_75sq_1325635414.jpg",
                     "username":"barackobama"
                  },
                  "id":"319634165414421617",
                  "text":"Thank you."
               },
               "comments":{
                  "count":28963,
                  "data":[
                     ...
                  ]
               },
               "created_time":"1352323376",
               "filter":"Nashville",
               "id":"319634107029710721_10206720",
               "images":{
                  "low_resolution":{
                     "height":306,
                     "url":"http://distilleryimage9.s3.amazonaws.com/4c4c3dc4292111e2af6f22000a1e9e28_6.jpg",
                     "width":306
                  },
                  "standard_resolution":{
                     "height":612,
                     "url":"http://distilleryimage9.s3.amazonaws.com/4c4c3dc4292111e2af6f22000a1e9e28_7.jpg",
                     "width":612
                  },
                  "thumbnail":{
                     "height":150,
                     "url":"http://distilleryimage9.s3.amazonaws.com/4c4c3dc4292111e2af6f22000a1e9e28_5.jpg",
                     "width":150
                  }
               },
               "imported":true,
               "likes":{
                  "count":258193,
                  "data":[
                     {
                        "full_name":"alyssataylor1",
                        "id":"337170532",
                        "profile_picture":"http://images.instagram.com/profiles/profile_337170532_75sq_1364472811.jpg",
                        "username":"alyssataylor1"
                     },
                    ...
                                      ]
               },
               "link":"http://instagram.com/p/RvkYbfmueB/",
               "location":null,
               "tags":[

               ],
               "type":"image",
               "user":{
                  "bio":"",
                  "full_name":"Barack Obama",
                  "id":"10206720",
                  "profile_picture":"http://images.instagram.com/profiles/profile_10206720_75sq_1325635414.jpg",
                  "username":"barackobama",
                  "website":""
               },
               "user_has_liked":false
            }
         },
        ...
      }
  },
  null
]
```
