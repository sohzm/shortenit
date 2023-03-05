# ShortenIt

URL Shortener using MySQL as DB.


</br>

### Getting started

You need to have [docker-compose](https://docs.docker.com/compose/) installed.

</br>

### Running it
```bash
docker-compose up
```
Now just open [localhost:9999](http://localhost:9999).             // btw you can configure everything like the port from the config file :)

</br>

### Working

Just running `docker-compose` will run 2 docker images, one will be the backend server and other one is mysql database. 

The backend is a go server which uses REST APIs for communicating with the frontend.

All the parameters like maximum links, expiry time for links, database credentials etc can be tweaked from the [config file](https://github.com/sz47/shortenit/blob/master/config.json).

</br>

### Requirements

* Storage of Max 20000 links, hence we only make alias of 3 characters.
* Links should expire in 24 hours.
* Need to use MySQL as DB.
