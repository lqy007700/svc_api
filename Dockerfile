FROM alpine
ADD svc_api /svc_api

ENTRYPOINT [ "/svc_api" ]