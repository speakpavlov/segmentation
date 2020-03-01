# Segmentation 

![Go](https://github.com/speakpavlov/segmentation/workflows/Go/badge.svg?branch=master)

curl -X GET \
  -H "Content-type: application/json" \
  -H "Accept: application/json" \
  -d '{"tag_id": "tag_1",\
         "data": {\
           "Level": 2,\
           "GPS": 1,\
           "UVS": 5\
         }}' \
  "http://localhost:9090/api/v1/segmentation"
  
  
ab -p post_loc.txt -T application/json -c 10 -n 100 http://localhost:9090/api/v1/segmentation
