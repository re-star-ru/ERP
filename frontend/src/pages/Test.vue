<template>
  <q-page padding>
    <h3 class="text-center">Опознание VIN кода</h3>
    <q-btn @click="navigatorMedia">распознать</q-btn>
    <video autoplay ref="cam"></video>
    <h5>
      <input
        @change="recognize"
        ref="imgCapture"
        type="file"
        accept="image/*"
      />
    </h5>
    <h4 class="text-center">
      Распознанный текст
      <br />
      {{ readyData }}
    </h4>
    <!--    <div class="row items-start justify-center">-->
    <!--      <SkuCard v-for="n in 20" :key="n"></SkuCard>-->
    <!--    </div>-->
  </q-page>
</template>

<script>
import SkuCard from 'components/SkuCard'
import Tesseract from 'tesseract.js'

export default {
  components: { SkuCard },
  data: () => {
    return {
      readyData: ''
    }
  },
  methods: {
    recognize() {
      console.log('regonize')
      let files = this.$refs.imgCapture.files
      console.log(files)
      Tesseract.recognize(
        files[0],
        // 'https://tesseract.projectnaptha.com/img/eng_bw.png',
        // this.$refs.cam,
        'rus+eng',
        {
          logger: m => {
            console.log(m)
            this.readyData = m.progress
          }
        }
      ).then(({ data: { text } }) => {
        console.log(text)
        this.readyData = text
      })
    },
    async navigatorMedia() {
      try {
        alert('vut')

        const stream = await navigator.mediaDevices.getUserMedia({
          video: {
            facingMode: { exact: 'environment' },
            width: { max: 600 },
            height: { max: 600 }
          }
          // video: true
        })
        this.$refs.cam.srcObject = stream
        alert(stream)
        this.recognize()
      } catch (e) {
        alert(e)
      }
    }
  }
}
</script>
