# Photobook

![upload page](https://github.com/fabrizio2210/photobook/blob/main/img/Screenshot_upload.png?raw=true)
![home page](https://github.com/fabrizio2210/photobook/blob/main/img/Screenshot_homepage.png?raw=true)

It's a simple Single Page Application that allows to upload photos and see them in the homepage.
I created it to have a place to share images for a specific event/party.  
The upload page is without access control, so everyone can upload photos.

## Features

- No login, I wanted to keep it simple as possible;
- The homepage is instantaneously updated when a new photo is uploaded;
- Appearance optimized for mobile, but also for desktop is fine (somewhat responsive).

## Environment Variables

`BLOCK_UPLOAD`: if set with a whatever value, it blocks the upload of photos.
`BLOCK_UPLOAD_MSG`: if `BLOCK_UPLOAD` is set, this is the custom message to show.

## Develop

### Requirements

To develop the webservice, you need:
- Docker
- docker-compose

### Steps to develop

```
user@host:~ cd photobook/
user@host:~/photobook docker/lib/createLocalDevStack.sh
```

Your development stack is reacheable at http://localhost/
You can modify the code from your GIT folder (so not in the Docker containers) and the servers (Vue and Flask) will automatically update.

## Deployment

The following steps will create a local stack in Docker. The script creates from the source the new images and start all the necessary components to serve the single page application.
```
user@host:~ cd photobook/
user@host:~/photobook docker/lib/createLocalStack.sh
```
Access it to http://localhost/

