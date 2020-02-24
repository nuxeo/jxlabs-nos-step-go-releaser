FROM scratch
EXPOSE 8080
ENTRYPOINT ["/step-go-releaser"]
COPY ./build/linux /