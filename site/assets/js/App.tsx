import React, { useState } from "react"
import { Filters } from "./components/Filters"

let arr: Item[] = []
if (process.env.NODE_ENV !== "production") {
  arr = [
    {
      amount: 1,
      char: "характеристика",
      id: "айди",
      name: "имя",
      type: "тип",
      sku: "артикул",
      images: [{ main: true, owner: "my", path: "kek" }],
    },
  ]
}

function App(): JSX.Element {
  const [text, setText] = useState<string>("")
  const [founded, setFounded] = useState<Item[]>(arr)
  const [loading, setLoading] = useState(false)
  const [searchAlert, setSearchAlert] = useState(false)

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setText(e.target.value)
  }

  function handleKey(e: React.KeyboardEvent) {
    if (e.key === "Enter") {
      searchHandler()
    }
  }

  function showSearchAlert() {
    setSearchAlert(true)
    setTimeout(() => {
      setSearchAlert(false)
    }, 3000)
  }

  async function searchHandler() {
    setLoading(true)

    try {
      const resp = await sendRequest(text)
      console.log("resp:", resp)
      setFounded(resp)
    } catch (e) {
      showSearchAlert()
      console.error("Error response")
    }

    setLoading(false)
  }

  return (
    <div>
      <div className='container mx-auto max-w-5xl p-4 flex '>
        <div className='flex-1 text-xl inline-block border-2 rounded-l-md border-yellow-400 shadow-xl'>
          <input
            className='w-full pl-2 py-2 rounded-md outline-none'
            type='text'
            name='lasdase'
            placeholder='Номер запчасти...'
            value={text}
            onChange={handleChange}
            onKeyDown={handleKey}
            id=''
          />
        </div>

        <button
          onClick={searchHandler}
          className='bg-yellow-400  pl-2 pr-8 py-2 rounded-l-none rounded-r-md shadow-xl cursor-pointer transition duration-300 ease-out hover:bg-yellow-500'
        >
          <svg
            role='status'
            className={`${
              loading ? "opacity-100" : "opacity-0"
            } inline align-text-bottom mr-2 w-5 h-5 text-gray-200 animate-spin dark:text-gray-600 fill-gray-600 dark:fill-gray-300`}
            viewBox='0 0 100 101'
            fill='none'
            xmlns='http://www.w3.org/2000/svg'
          >
            <path
              d='M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z'
              fill='currentColor'
            />
            <path
              d='M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z'
              fill='currentFill'
            />
          </svg>
          Найти
        </button>
      </div>

      {searchAlert && (
        <div className='absolute inset-x-0 mx-auto max-w-xl'>
          <SearchAlert />
        </div>
      )}

      <div className='container  mx-auto max-w-3xl px-4 pb-2 border-b  border-gray-200'>
        <Filters />
      </div>
      <div>
        <CardList cards={founded} />
      </div>
    </div>
  )
}
export { App }

function SearchAlert() {
  return (
    <div
      className='flex p-4 mb-4 text-sm text-yellow-700 bg-yellow-100 rounded-lg dark:bg-yellow-200 dark:text-yellow-800'
      role='alert'
    >
      <svg
        className='inline flex-shrink-0 mr-3 w-5 h-5'
        fill='currentColor'
        viewBox='0 0 20 20'
        xmlns='http://www.w3.org/2000/svg'
      >
        <path
          fillRule='evenodd'
          d='M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z'
          clipRule='evenodd'
        ></path>
      </svg>
      <div>
        <span className='font-medium'>Ошибка на сервере!</span>
      </div>
    </div>
  )
}

// images: Array(5)
// 0:
// main: true
// owner: "77367274-a5f1-41b7-80f4-d9e6f0070840"
// path: "srv1c/images/48fceb63-ae56-4f35-9635-8514b97763b7.jpeg"
// [[Prototype]]: Object
// 1: {main: false, owner: '77367274-a5f1-41b7-80f4-d9e6f0070840', path: 'srv1c/images/fbb83137-4b51-4b21-9900-e26c620272fa.jpeg'}
// 2: {main: false, owner: '77367274-a5f1-41b7-80f4-d9e6f0070840', path: 'srv1c/images/14efc2f0-9f68-484c-8d39-a5414e23bf90.jpeg'}
// 3: {main: false, owner: '77367274-a5f1-41b7-80f4-d9e6f0070840', path: 'srv1c/images/c0d048cf-8ccd-4c91-8373-9f3b263e43e2.jpeg'}
// 4: {main: false, owner: '77367274-a5f1-41b7-80f4-d9e6f0070840', path: 'srv1c/images/6822ef7b-81d6-40df-8b24-1e7edaf398ee.jpeg'}
// length: 5
// [[Prototype]]: Array(0)

let url = "https://api.re-star.ru/v1/oprox"
let noImage = "https://via.placeholder.com/800x600"

if (process.env.NODE_ENV !== "production") {
  url = "http://localhost:8100"
  noImage = "https://loremflickr.com/800/600/cat"
}

type Image = {
  main: boolean
  owner: string
  path: string
}

type Item = {
  amount: number
  type: string
  char: string
  id: string
  name: string
  sku: string
  images: Image[]
}

type CardListProps = {
  cards: Item[]
}

type ImageListProps = {
  images: Image[]
}

function s3path(path: string): string {
  return `https://s3.re-star.ru/${path}`
}

function ImageList(props: ImageListProps) {
  const images = props.images

  const imageList = images.map((image) => (
    <li key={image.path}>
      <img src={s3path(image.path)} alt={image.owner} />
    </li>
  ))

  return <ul>{imageList}</ul>
}

function CardList(props: CardListProps) {
  const cards = props.cards

  const listCards = cards.map((card) => (
    <li
      className='max-w-sm bg-white rounded-lg shadow-md dark:bg-gray-800 dark:border-gray-700 p-4'
      key={card.id}
    >
      <ImageList images={card.images} />

      <a href={`/item/${card.name}`}>
        <img className='rounded-t-lg' src={noImage} alt='product image' />
      </a>

      <h3 className='pt-4 text-xl font-semibold tracking-tight text-gray-900 dark:text-white '>
        {card.type} {card.name}, {card.char}
      </h3>
      <h4 className='pt-2 font-semibold  text-gray-900 dark:text-white '>
        Арт. {card.sku}
      </h4>
      <h4 className='text-md font-bold text-gray-900 dark:text-white'>
        В наличии {card.amount} шт.
      </h4>
    </li>
  ))

  return (
    <ul className='p-4  flex flex-wrap gap-4 justify-center'>{listCards}</ul>
  )
}

async function sendRequest(r: string): Promise<Item[]> {
  let response = await fetch(url + /search/ + r)
  if (response.ok) {
    let json = await response.json()
    return json
  } else {
    alert("ошибка HTTP: " + response.status)
    return []
  }
}
