watch-css:
	npx tailwindcss -i ./styles.css -o ./public/styles.css --watch

dev:
	air -c ./.air.toml & \
	npx tailwind \
		-i 'styles.css' \
		-o 'public/styles.css' \
		--watch