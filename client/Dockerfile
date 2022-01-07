FROM quay.io/jitesoft/node
#FROM openshift/nodejs:latest
WORKDIR /app
COPY package.json ./
COPY package-lock.json ./
COPY ./ ./
#RUN npm ci && npm build && rm -rf node_modules
RUN npm i
EXPOSE 8080
RUN chown -R 1001:0 /app
USER 1001
CMD ["npm", "run", "start"]

#FROM node:14
#FROM registry.access.redhat.com/ubi8/nodejs-14
#USER root
##WORKDIR /usr/src/app
#WORKDIR /tmp/src
##ENV NODE_ENV production
#
## Just copy the package.json...
#COPY package*.json ./
#
##COPY ./ ./
#COPY build/ ./
##RUN npm i
#
## so we can cache this layer:
##RUN npm install
#
#EXPOSE 8080
##RUN chown -R 1001:0 /usr/src/app
#USER 1001
## command to run on container start
#CMD [ "npm", "run", "start"]

# Stage 1: Use yarn to build the app
#FROM node:14 as builder
#WORKDIR /usr/src/app
#COPY package.json ./
#RUN yarn
#COPY . ./
#RUN yarn build

#FROM registry.access.redhat.com/ubi8/nodejs-14 as builder
#USER root
#WORKDIR /usr/src/app
#COPY package.json ./
#COPY ./ ./
#RUN npm install
##COPY . ./
##RUN yarn build
#
## Stage 2: Copy the JS React SPA into the Nginx HTML directory
#FROM bitnami/nginx:latest
#COPY --from=builder /usr/src/app/build /app
##COPY . ./
#EXPOSE 8080
##USER 1001
#CMD ["nginx", "-g", "daemon off;"]