if ($args[0] -eq "test") {
  go test --cover --race -v $(glide novendor)
} elseif ($args[0] -eq "bootstrap") {
  glide install
} else {
  echo "Unknown command: '$args'"
  exit 1
}

exit $LASTEXITCODE
