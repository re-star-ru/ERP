import React, { useState } from "react" // import * as React from "react"

// import Card from "./components/Card"

function App() {
  const [text, setText] = useState<string>("")
  const [founded, setFounded] = useState<Item[]>([
    { amount: 1, char: "характеристика", id: "айди", name: "имя", type: "тип" },
  ])

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setText(e.target.value)
  }

  function handleKey(e: React.KeyboardEvent) {
    if (e.key === "Enter") {
      searchHandler()
    }
  }

  async function searchHandler() {
    const resp = await sendRequest(text)
    console.log("resp:", resp)
    setFounded(resp)
  }

  return (
    <div>
      <div className='container mx-auto p-4 flex '>
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
        <input
          className='bg-yellow-400 px-4 py-2 rounded-r-md shadow-xl cursor-pointer transition duration-300 ease-out hover:bg-yellow-500'
          type='button'
          value='Найти'
          onClick={searchHandler}
        />
      </div>

      <div>
        <CardList cards={founded} />
      </div>
    </div>
  )
}
export { App }

// const url = "https://api.re-star.ru/v1/oprox"
const url = "http://localhost:8100"

type Item = {
  amount: number
  type: string
  char: string
  id: string
  name: string
}

type CardListProps = {
  cards: Item[]
}

function CardList(props: CardListProps) {
  const cards = props.cards
  const listCards = cards.map((card) => (
    <li
      className='max-w-sm bg-white rounded-lg shadow-md dark:bg-gray-800 dark:border-gray-700 p-6'
      key={card.id}
    >
      <a href={`/item/${card.name}`}>
        <img
          className='rounded-t-lg'
          src='https://loremflickr.com/800/600/cat'
          alt='product image'
        />
      </a>

      <h3 className='mt-4 text-xl font-semibold tracking-tight text-gray-900 '>
        {card.type} {card.name}, {card.char}
      </h3>
      <h4 className='mt-4'>В наличии {card.amount} шт.</h4>
    </li>
  ))

  return <ul className='p-10  flex flex-wrap gap-4'>{listCards}</ul>
}

async function sendRequest(r: string): Promise<Item[]> {
  let response = await fetch(url + /search/ + r)
  if (response.ok) {
    let json = await response.json()
    // console.log("resp", json)
    return json
  } else {
    alert("ошибка HTTP: " + response.status)
    return []
  }
}
