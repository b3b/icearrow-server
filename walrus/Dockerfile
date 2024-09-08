FROM mysten/sui-tools:testnet

ARG USER_ID=1000
ARG GROUP_ID=1000

RUN  \
     apt-get update \
     && apt-get install -y curl \
     && rm -rf /var/lib/apt/lists/*

RUN curl -o /usr/local/bin/walrus \
         https://storage.googleapis.com/mysten-walrus-binaries/walrus-latest-ubuntu-x86_64 \
    && chmod +x /usr/local/bin/walrus

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN groupadd -g ${GROUP_ID} sui \
    && useradd -m -u ${USER_ID} -g ${GROUP_ID} sui \
    && mkdir -p /home/sui/.sui/sui_config \
    && chown sui:sui /home/sui/.sui/sui_config
VOLUME /home/sui/.sui/sui_config
USER sui

RUN \
    mkdir -p ~/.config/walrus \
    && curl -o ~/.config/walrus/client_config.yaml \
       https://storage.googleapis.com/mysten-walrus-binaries/walrus-configs/client_config.yaml


ENTRYPOINT ["/entrypoint.sh"]
CMD ["walrus", "daemon", "-b", "0.0.0.0:31415"]
