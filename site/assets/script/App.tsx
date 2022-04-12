import React, { useState } from "react" // import * as React from "react"

// import Card from "./components/Card"

function App() {
  const [text, setText] = useState("")

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    console.log(e.target.value)
    setText(e.target.value)
  }

  function searchHandler() {
    const resp = sendRequest(text)
    console.log("resp:", resp)
    setText("")
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
    </div>
  )
}
export { App }

function sendRequest(r: string): string {
  return "ok"
}
