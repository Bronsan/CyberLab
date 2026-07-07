#!/bin/bash
# Generate dynamic flag on container startup
FLAG="flag{$(cat /proc/sys/kernel/random/uuid 2>/dev/null || uuidgen 2>/dev/null || echo "debug-$(date +%s)")}"
echo "$FLAG" > /flag.txt
chmod 444 /flag.txt
echo "Dynamic flag generated: $FLAG"

# Start Apache
exec apache2-foreground
