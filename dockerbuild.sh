docker image build -f Dockerfile -t forumtask .

docker system prune

docker images

docker container run -it --rm -p 8888:8888 --name forum-docker forumtask

docker ps -a