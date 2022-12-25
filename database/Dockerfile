# Use this image
FROM postgres:10.3

# Copy up.sql file into speicified directory
COPY up.sql /docker-entrypoint-initdb.d/1.sql

# run command
CMD ["postgres"]