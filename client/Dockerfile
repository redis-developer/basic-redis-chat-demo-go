FROM node:alpine
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
COPY ./ ./
RUN npm i
EXPOSE 8080
CMD ["npm", "run", "start"]

##FROM node:14
#FROM registry.access.redhat.com/ubi8/nodejs-14-minimal
#
#WORKDIR /usr/src/app
#
##ENV NODE_ENV production
#
## Just copy the package.json...
#COPY package*.json ./
#
#COPY ./ ./
#RUN npm i
#
## so we can cache this layer:
##RUN npm install
#
#EXPOSE 8080
#
## command to run on container start
#CMD [ "npm", "run", "start"]