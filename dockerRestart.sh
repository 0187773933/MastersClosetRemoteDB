#!/bin/bash
APP_NAME="public-mct-remote-db"
id=$(sudo docker restart $APP_NAME)
sudo docker logs -f $id