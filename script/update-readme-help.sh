#!/usr/bin/env sh

set -e -o pipefail

cd "$(dirname "$0")/.."

README="README.md"

HELP=$(go run . --help 2>&1)

awk -v "helptext=$HELP" '
    BEGIN { in_block=0 }
    /<!-- HELP_START -->/ {
        print;
        print "<!-- DO NOT EDIT BELOW THIS LINE! This section is auto-generated. -->";
        print "";
        print "```";
        print "$ bloom --help";
        print helptext;
        print "```";
        print "";
        in_block=1;
        next
    }
    /<!-- HELP_END -->/ { in_block=0 }
    !in_block
' "$README" >"${README}.tmp" && mv "${README}.tmp" "$README"

echo "✓ updated $README with latest help output"
