#/bin/bash

 docker ps | grep api:dev | awk '{print $1}' | xargs docker logs -f
 