./target/container -api https://api/ -dir .data/ai-docker/\
  -github wass88/reversi-random\
  -branch master\
  -commit b3cd1a475dded156758005866761de51ee690607\
  setup

./target/container -api https://api/ -dir .data/ai-docker/\
  -github wass88/reversi-random\
  -branch master\
  -commit b3cd1a475dded156758005866761de51ee690607\
  -cpu 100 -mem 256\
  run