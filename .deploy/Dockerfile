FROM node:20-alpine as build

COPY package.json .
COPY yarn.lock .

RUN yarn install --frozen-lockfile

COPY . .

RUN yarn build


FROM node:20-alpine as install-node-modules

COPY --from=build /package.json ./package.json
COPY --from=build /yarn.lock    ./yarn.lock

RUN yarn install --production --frozen-lockfile


FROM node:20-alpine as deploy

WORKDIR /opt/app

COPY --from=build                   /dist          ./dist/
COPY --from=install-node-modules    /node_modules  ./node_modules

ENV NODE_ENV=production

CMD ["node", "dist/main"]
