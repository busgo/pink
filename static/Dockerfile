FROM node:16.3.0

RUN \
npm config set registry https://registry.npm.taobao.org \
&& npm install -g @vue/cli
WORKDIR /app
COPY .  ./
RUN npm install cnpm -g
RUN cnpm install
EXPOSE 8002
ENTRYPOINT cnpm run dev
