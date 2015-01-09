#!/bin/bash

if [ -d "$FRONT_END_DIR" ]; then
	mkdir -p collect/api/front_end/
	cp -r ${FRONT_END_DIR}/* collect/api/front_end/
else
	echo "Front end directory does not exist, ignoring: $FRONT_END_DIR"
fi
