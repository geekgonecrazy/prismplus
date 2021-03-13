# prism+

Based on prism but with api that allows you to dynamically add new sessions with session keys and dynamically add and remove destinations.


## Create Session
```
curl -L -X POST 'http://localhost:5383/api/v1/sessions' \
-H 'Content-Type: application/json' \
--data-raw '{
    "key": "abc1234"
}'
```

You can now start streaming to localhost:1935/live/abc1234

## Add Destination
```
curl -L -X POST 'http://localhost:5383/api/v1/sessions/abc1234/destinations' \
-H 'Content-Type: application/json' \
--data-raw '{
    "url": "rtmp://broadcast.owncast.online/live/octempdemoazfdhw"
}'
```

You should then start to see the content streaming there

## Remove Destination
```
curl -L -X DELETE 'http://localhost:5383/api/v1/sessions/abc1234/destinations/0'
```

It will stop playing at that destination

## End Session
```
curl -L -X DELETE 'http://localhost:5383/api/v1/sessions/abc1234'
```


