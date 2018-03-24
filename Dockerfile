FROM scratch

COPY tool-cal tool-cal

EXPOSE 8080

CMD [ "./tool-cal" ]