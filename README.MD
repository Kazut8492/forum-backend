## How to use
1. Install Docker to your pc
2. Within command line window, nagivate to the project folder
3. Run `bash dockerbuild.sh` and go to `localhost:8888`.
### Or you can also do it manually as well...
3. Run `docker image build -f Dockerfile -t forum .` to create docker image
4. Run `docker container run -it --rm -p 8888:8888 --name forum-docker forum` to create and run docker container
  ** In case that the system says that the container is already in use, then Run `docker rm <docker container name>` like `docker rm forum-docker`
5. Type `localhost:8888` in the browser to jump to the webpage
### ---For Audit---
1. Open a new command line window, nagivate to the project folder again
2. Run `docker images` to get a list of docker images
3. Run `docker ps -a` to get a list of docker containers

## Authors
* Kazuki
* Semjon
* SlabiJ