FROM php:8.2-apache

# Install mysqli
RUN docker-php-ext-install mysqli

# Copy challenge source
COPY src/ /var/www/html/

# Generate dynamic flag at container start
COPY flag/generate.sh /generate.sh
RUN chmod +x /generate.sh

# Set up entrypoint to generate flag on start
COPY start.sh /start.sh
RUN chmod +x /start.sh

ENTRYPOINT ["/start.sh"]
CMD ["apache2-foreground"]
