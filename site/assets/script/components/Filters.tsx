import React, { useState } from "react"
import { Transition } from "@headlessui/react"

function Filters(): JSX.Element {
  const [inStock, setInStock] = useState(true)

  function stockHandler() {
    setInStock(!inStock)
    console.log("in sotk" + inStock)
  }

  const [isShowing, setIsShowing] = useState(false)

  return (
    <>
      <div className='inline-block p-1 w-36 rounded-lg bg-white border'>
        <span className='cursor-pointer' onClick={stockHandler}>
          В наличии
        </span>
        <Transition show={true} as='span'>
          ✅
        </Transition>
      </div>

      <button onClick={() => setIsShowing((isShowing) => !isShowing)}>
        Toggle
      </button>
      <Transition show={isShowing}>I will appear and disappear.</Transition>
    </>
  )
}

export { Filters }
