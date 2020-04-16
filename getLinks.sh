#! /bin/bash
echo "Enter the query to search"
read query
IFS=' ' read -r params <<<$query
function join {
	local IFS="$1"
	shift
	echo "$*"
}

query=$(join + ${params[@]})
echo "Don't freak out, now I'm gonna download the respective html, to extract the real links."
sleep 4
curl -X GET -L "https://xnxx.com/search/$query">response.html
grep 'href="' response.html | cut -d '=' -f7 | cut -d '"' -f2 | grep "/video-">list
base="https://xnxx.com"
for link in $(cat list);
do
    link="$base$link"
    echo "entering in $link"
    curl -X GET -L $link>all.html
    url="$(grep "setVideoUrlHigh" all.html | cut -d "'" -f2)&amp;download=1"
    echo $url>>realLinks
done
echo "Cleaning files"
rm list response.html all.html
echo "The links are in realLinks, in the current directory"
echo "Use curl or wget to download each of this links"
echo "Good luck :), perverso"
