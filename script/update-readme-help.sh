#!/usr/bin/env sh

set -e

cd "$(dirname "$0")/.."

README="readme.md"
BINARY="./rose-pine-bloom"
HELP_TMP="$(mktemp)"

if [ ! -x "$BINARY" ]; then
	echo "Building $BINARY..."
	go build -o $BINARY .
fi

echo "\$ $BINARY --help" >"$HELP_TMP"
$BINARY --help >>"$HELP_TMP" 2>&1

awk -v helpfile="$HELP_TMP" '
    BEGIN { in_block=0 }
    /<!-- HELP_START -->/ {
        print;
        print "<!-- DO NOT EDIT BELOW THIS LINE! This section is auto-generated. -->";
        print "";
        print "```";
        while ((getline line < helpfile) > 0) print line;
        close(helpfile);
        print "```";
        in_block=1;
        next
    }
    /<!-- HELP_END -->/ { in_block=0 }
    !in_block
' "$README" >"${README}.tmp" && mv "${README}.tmp" "$README"

rm "$HELP_TMP"

echo "✓ Updated $README with latest help output."
