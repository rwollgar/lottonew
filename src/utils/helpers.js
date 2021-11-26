module.exports = {

    env: {
        theme:{
            gamepanelbg:'#E6F4F1',
            bordercolor:'#AA5EA2',
            numberbg:'#7FD1AE'
        }
    },    

    clock: {
        sleep: (milliseconds) => { return new Promise(resolve => setTimeout(resolve, milliseconds)) }
    }
}