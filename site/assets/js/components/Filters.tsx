import React, { useState } from "react"

function Filters(): JSX.Element {
  const [inStock, setInStock] = useState(true)
  function stockHandler() {
    setInStock(!inStock)
  }

  return (
    <>
      <div className='inline-block p-1 rounded-lg bg-white border'>
        <span className='cursor-pointer' onClick={stockHandler}>
          В наличии
          <span
            className={`transition-opacity  ${
              inStock ? "opacity-100" : "opacity-0"
            }`}
          >
            ✅
          </span>
        </span>
      </div>
    </>
  )
}

export { Filters }
