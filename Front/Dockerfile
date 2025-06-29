

FROM node:20 AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .

ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=$REACT_APP_API_URL

RUN npm run build


FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]


# docker build . -t kostinalex/jule_front 
# docker run -p 49000:80 --name jule_front -d kostinalex/jule_front
# docker push kostinalex/jule_front
# docker pull kostinalex/jule_front
