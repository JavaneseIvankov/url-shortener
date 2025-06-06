#!/usr/bin/env bash

BASE_URL="http://localhost:4321"
API_URL="$BASE_URL/api/v1"

err() {
  echo "$*" >&2
}

out() {
  echo "$*" >&1
}

saveAccessToken() {
  echo "$1" > ./tmp/access-token.txt
}

saveRefreshToken() {
  echo "$1" > ./tmp/refresh-token.txt
}

getAcessToken() {
  cat ./tmp/access-token.txt
}

getRefreshToken() {
  cat ./tmp/refresh-token.txt
}

displayRespose() {
  res=$1
  echo "==== Reponse ==="
  echo "$1"
  echo "==== end ==="
}

register() {
  TEST_EMAIL='janedoe@gmail.com'
  TEST_PASSWORD='password'

  res=$(http POST "$API_URL/auth/register/" email="$TEST_EMAIL" password="$TEST_PASSWORD")
  accessToken=$(echo "$res" | jq -r '.access_token')
  refreshToken=$(echo "$res" | jq -r '.refresh_token')
  displayRespose "$res"
  if [ -z "$accessToken" ] || [ "$accessToken" == "null" ] || [ -z "$refreshToken" ] || [ "$refreshToken" == "null" ]; then
      return 1
  else
      saveAccessToken "$accessToken"
      saveRefreshToken "$refreshToken"
      return 0
  fi
}

login() {
  TEST_EMAIL='janedoe@gmail.com'
  TEST_PASSWORD='password'

  res=$(http POST "$API_URL/auth/login/" email="$TEST_EMAIL" password="$TEST_PASSWORD")
  accessToken=$(echo "$res" | jq -r '.access_token')
  refreshToken=$(echo "$res" | jq -r '.refresh_token')
  displayRespose "$res"
  if [ -z "$accessToken" ] || [ "$accessToken" == "null" ] || [ -z "$refreshToken" ] || [ "$refreshToken" == "null" ]; then
      return 1
  else
      saveAccessToken "$accessToken"
      saveRefreshToken "$refreshToken"
      return 0
  fi
}

shortenUrl() {
  echo "Running shorten..."
  SHORTEN_URL="https://google.com"
  SHORT_NAME="ggl"

  accessToken=$(getAcessToken)
  res=$(http POST "$API_URL/shorten" Authorization:"Bearer $accessToken" url=$SHORTEN_URL short_name=$SHORT_NAME)
  displayRespose "$res"
  if [ -z "$res" ] || [ "$res" == "null" ]; then
      return 1
  else
      return 0
  fi
}

editUrl() {
  echo "Running edit url..."
  SHORT_NAME="ggl"
  NEW_URL="https://youtube.com"

  accessToken=$(getAcessToken)
  res=$(http POST "$API_URL/edit/$SHORT_NAME" Authorization:"Bearer $accessToken" new_url=$NEW_URL)
  displayRespose "$res"
  if [ -z "$res" ] || [ "$res" == "null" ]; then
      return 1
  else
      return 0
  fi
}

refresh() {
  out "Running refresh..."

  token=$(getRefreshToken)
  res=$(http POST "$API_URL/auth/refresh-session/" refresh_token="$token")
  accessToken=$(echo "$res" | jq -r '.access_token')
  refreshToken=$(echo "$res" | jq -r '.refresh_token')
  displayRespose "$res"
  if [ -z "$accessToken" ] || [ "$accessToken" == "null" ] || [ -z "$refreshToken" ] || [ "$refreshToken" == "null" ]; then
      return 1
  else
      saveAccessToken "$accessToken"
      saveRefreshToken "$refreshToken"
      return 0
  fi
}

case "$1" in 
  register)
    register
    if [ $? -eq 0 ]; then
      out "Sucessfully registered"
    else 
      err "Something went wrong"
    fi
    ;;
  shorten)
    shortenUrl
    if [ $? -eq 0 ]; then
      out "Successfuly shortened url"
    else 
      err "Shorten url request failed"
    fi
    ;;
  edit)
    editUrl
    if [ $? -eq 0 ]; then
      out "Successfuly edit short url"
    else 
      err "Edit short url request failed"
    fi
    ;;
  login)
    login
    if [ $? -eq 0 ]; then
      out "Successfuly login"
    else 
      err "Login failed"
    fi
    ;;
  refresh)
    refresh
    if [ $? -eq 0 ]; then
      out "Session refreshed Successfuly"
    else 
      err "Refresh session attempt failed"
    fi
    ;;
  *)
    err "Invalid option selected"
    ;;
esac
