ALB_URL="https://nlp.example-api.com"
TEXT="Linux, computer operating system created in the early 1990s by Finnish software engineer Linus Torvalds and the Free Software Foundation (FSF). While still a student at the University of Helsinki, Torvalds started developing Linux to create a system similar to MINIX, a UNIX operating system."
TEXT_ALT1="Today, the Nobel Prizes are regarded as the most prestigious awards in the world in their various fields. Notable winners have included Marie Curie, Theodore Roosevelt, Albert Einstein, George Bernard Shaw, Winston Churchill, Ernest Hemingway, Martin Luther King, Jr., the Dalai Lama, Mikhail Gorbachev, Nelson Mandela and Barack Obama. Multiple leaders and organizations sometimes receive the Nobel Peace Prize, and multiple researchers often share the scientific awards for their joint discoveries. In 1968, a Nobel Memorial Prize in Economic Science was established by the Swedish national bank, Sveriges Riksbank, and first awarded in 1969."
API_KEY="DqiSyCzJgUY9kbxWiF1QA7NY"

for i in {1..10}
do
  ROUTE=health
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X GET \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" | jq

  ROUTE=health/lang
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X GET \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" | jq

  ROUTE=health/prose
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X GET \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" | jq

  ROUTE=health/rake
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X GET \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" | jq

  ROUTE=health/dynamo
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X GET \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" | jq

  ROUTE=routes
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X GET \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" | jq

  ROUTE=keywords
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X POST \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"text\": \"${TEXT_ALT1}\"}" | jq

  ROUTE=language
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X POST \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"text\": \"${TEXT_ALT1}\"}" | jq

  ROUTE=entities
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X POST \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"text\": \"${TEXT_ALT1}\"}" | jq

  ROUTE=tokens
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X POST \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"text\": \"${TEXT_ALT1}\"}" | jq

  ROUTE=sentences
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X POST \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"text\": \"${TEXT_ALT1}\"}" | jq

  ROUTE=record
  printf "\n/%s\n" "${ROUTE}"
  curl -s -X POST \
    ${ALB_URL}/${ROUTE} \
    -H "Content-Type: application/json" \
    -H "X-API-Key: ${API_KEY}" \
    -d "{\"text\": \"${TEXT_ALT1}\"}" | jq

###### THESE THROW ERRORS ######

  # ROUTE=error
  # printf "\n/%s\n" "${ROUTE}"
  # curl -s -X GET \
  #   ${ALB_URL}/${ROUTE} \
  #   -H "Content-Type: application/json" \
  #   -H "X-API-Key: ${API_KEY}" | jq

  # # Special cases
  # # No Auth
  # ROUTE=language
  # printf "\n/%s (No Auth - 403 Forbidden)\n"  "${ROUTE}"
  # curl -s -X POST \
  #   ${ALB_URL}/${ROUTE} \
  #   -H "Content-Type: application/json" \
  #   -d "{\"text\": \"SELECT * FROM passwords;\"}"

  # # Bad Endpoint
  # ROUTE=bad-endpoint
  # printf "\n/%s (Not Found)\n"  "${ROUTE}"
  # curl -s -X GET \
  #   ${ALB_URL}/${ROUTE} \
  #   -H "Content-Type: application/json" \
  #   -H "X-API-Key: ${API_KEY}" | jq

  # # Should be blocked by WAF:
  # # AWS#AWSManagedRulesSQLiRuleSet#SQLi_BODY
  # ROUTE=language
  # printf "\n/%s (WAF Rule - SQLi_BODY)\n" "${ROUTE}"
  # curl -s -X POST \
  #   ${ALB_URL}/${ROUTE} \
  #   -H "Content-Type: application/json" \
  #   -H "X-API-Key: ${API_KEY}" \
  #   -d "{\"text\": \"SELECT * FROM passwords;\"}"

  # # Should be blocked by WAF:
  # # AWS#AWSManagedRulesCommonRuleSet#CrossSiteScripting_QUERYARGUMENTS
  # ROUTE=language
  # printf "\n/%s (WAF Rule - CrossSiteScripting_QUERYARGUMENTS)\n" "${ROUTE}"
  # curl -s -X POST \
  #   "${ALB_URL}/${ROUTE}?name=Bob%0d%0a%0d%0a<script>alert(document.domain)</script>" \
  #   -H "Content-Type: application/json" \
  #   -H "X-API-Key: ${API_KEY}" \
  #   -d "{\"text\": \"WAF Cross Site Scripting test - body\"}"
done
