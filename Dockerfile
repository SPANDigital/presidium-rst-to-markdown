FROM scratch
ENTRYPOINT ["/rst2md"]
COPY dist/presidium-rst-to-markdown_linux_amd64_v1/rst2md /