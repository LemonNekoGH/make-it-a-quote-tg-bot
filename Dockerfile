FROM node:16-alpine

WORKDIR /app/

COPY . .

RUN npm i -g pnpm
# 安装 node-canvas 需要的包
RUN apk add pixman-dev cairo-dev pango-dev
# 安装 alpine linux 缺少的包
RUN apk add python3 make g++ pkgconfig
RUN pnpm i

CMD ["pnpm", "dev"]