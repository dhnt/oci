FROM alpine:edge

# add openssh and clean
RUN apk add --update openssh \
&& rm  -rf /tmp/* /var/cache/apk/*
# add entrypoint script
ADD sshd-entrypoint.sh /usr/local/bin
#make sure we get fresh keys
RUN rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key

EXPOSE 22

ENV USER=app
RUN adduser -D -h /$USER -u 1000 -G users $USER \
    && echo $USER:$USER | chpasswd

ENTRYPOINT ["sshd-entrypoint.sh"]
CMD ["/usr/sbin/sshd","-D"]
