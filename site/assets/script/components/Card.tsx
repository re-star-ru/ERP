import { defineComponent, ref, h } from "vue"

export default defineComponent({
  setup() {
    const count = ref(0)
    const increment = () => count.value++
    return {
      count,
      increment,
    }
  },

  render() {
    return (
      <div>
        <div class='container mx-auto p-4 flex '>
          <div class='flex-1 text-xl inline-block border-2 rounded-l-md border-yellow-400 shadow-xl'>
            <input
              class='w-full pl-2 py-2 rounded-md outline-none'
              type='text'
              name='lasdase'
              placeholder='Номер запчасти...'
              value={this.count}
              id=''
            />
          </div>
          <input
            class='bg-yellow-400 px-4 py-2 rounded-r-md shadow-xl cursor-pointer transition duration-300 ease-out hover:bg-yellow-500'
            type='button'
            value='Найти'
            onClick={this.increment}
          />
        </div>
      </div>
    )
  },
})
