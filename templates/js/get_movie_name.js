let fetchDetails = async () => {
  const response = await fetch("/list-movies")
  const data = await response.json()
  return data
}


async function getDetails() {
  return await fetchDetails()

}