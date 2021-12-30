import create from 'zustand';

import produce from 'immer';
import pipe from 'ramda/es/pipe';
//import fetch from 'node-fetch';
import fetch from 'isomorphic-unfetch'

import lotto1 from '../assets/img/lottoballs_1.jpg';
import lotto2 from '../assets/img/lottoballs_2.jpg';
import lotto3 from '../assets/img/lottoballs_3.jpg';

const apiUrl = 'http://localhost:1337/api';

// const myhook_mw = config => (set, get, api) => config(args => {
//     console.log('BEFORE MYHOOK: ', get().selectedDraw);
//     set(args);
//     console.log('AFTER MYHOOK: ', get().selectedDraw);
// }, get, api);

const immer_mw = config => (set, get, api) => config(fn => {
    //console.log('BEFORE: ', get().selectedDraw);
    set(produce(fn))
    //console.log('AFTER: ', get().selectedDraw);
}, get, api);

const log_mw = config => (set, get, api) => config(args => {
    console.log("  applying", args)
    set(args)
    //console.log("  new state", get())
}, get, api);

//const createStore = pipe(log_mw, myhook_mw, immer_mw, create);
const createStore = pipe(log_mw, immer_mw, create);

const useLnStore = createStore((set) => ({

    initialised: false,
    refreshing: false,
    lastError: null,
    data: [],
    gameData: null,
    selectedDraw: null,
    imageMap: {
        'oz-lotto':lotto1,
        'powerball':lotto2,
        'monday-lotto':lotto3,
        'wednesday-lotto':lotto3,
        'saturday-lotto':lotto3
    },
    
    // setProcessData: (value) => set((state) => {state.processData = value}),
    // setRefreshing: (value) => set((state) => {state.refreshing = value}),
    // setCurrentProcessId: (value) => set((state) => {state.currentProcessId = value}),
    // setCurrentTab: (pid, tab) => set((state) => {state.currentTab[pid] = tab}),
    //setCurrentProcessStatus: (pid, value) => set((state.processStatusNew[pid] = value)),
    // setTimezone: () => set((state) => {state.currentTZ = state.currentTZ === 'UTC' ? 'LOCAL' : 'UTC'}),

    initStore: () => set(async (state) => {

        console.log('initStore...');

        try {
            // const pInfo = await fetch(`${apiUrl}/games`, {mode: 'cors'});
            // const data = await pInfo.json();
            const data = await dataApi.getGames();
            set(() => ({
                initialised: true, 
                data: data.sort((a,b) => {return a.order - b.order;})
            }))
        } catch (error) {
            console.log(error);
        }
    }),

    getGameData: (game) => set(async (state) => {

        // const uri = `${rsUrl}/games/${game}`;
        // console.log('Get Game: ', uri);
        // const data = await fetch(uri, { mode: 'cors', method: 'GET' });
        const data = await dataApi.getGameData(game);
        set(() => ({ gameData: data }));
        
    }),

    setSelectedDraw: (id) => set((state) => {
        //console.log('setSelectedDraw: ', id, state.selectedDraw);
        state.selectedDraw = id;
        //set((id) => ({ selectedDraw: id }));
    })

    // startRefresh: (value) => set((state) => {
        
    //     console.log('zustandStore:startRefresh: ', value);

    //     setInterval(async (rsUrl) => {
    
    //         set(() => {return {refreshing:true}});
            
    //         try {
    //             const pInfo = await fetch(`${apiUrl}/processinfo`, {mode: 'cors'});
    //             const data = await pInfo.json();
    //             set(() => ({interval: value, processData: data, refreshing: false}))
    //         } catch (error) {
    //             set(() => ({lastError: error, refreshing: false}));
    //             console.log(error)
    //         }
    
    //     }, value, rsUrl);

    // }),

    // stopProcess: (value) => set(async (state) => {

    //     //set((state.setCurrentProcessStatus[value] = 'stopping...'));
    //     //set(() => ({processData: data, processStatus: 'stopping...'}));
        
    //     console.log('stopProcess: ', value);
    //     let uri = `${apiUrl}/processes/stop/${value}`;
    //     let result = await fetch(uri, {mode: 'cors', method: 'GET'});
    //     result = await result.json;
    
    //     if(result) {
    //         console.log(result);
    //     }
       
    //     const pInfo = await fetch(`${apiUrl}/processinfo`, {mode: 'cors'});
    //     const data = await pInfo.json();
    //     set(() => ({processData: data})) //, processStatus: 'stopped'}));
    //     //set((state.setCurrentProcessStatus[value] = 'stopped'));

    // }),

    // startProcess: (value) => set(async (state) => {

    //     //set(() => ({processStatus: 'starting...'}));
    //     console.log('startProcess: ', value);
    //     let uri = `${rsUrl}/processes/start/${value}`;
    //     let result = await fetch(uri, {mode: 'cors', method: 'GET'});

    //     if(result) {
    //         result = await result.json;
    //         console.log(result);
    //     }

    //     setTimeout(async () => {
    //         const pInfo = await fetch(`${rsUrl}/processinfo`, {mode: 'cors'});
    //         const data = await pInfo.json();
    //         set(() => ({processData: data})); //, processStatus:'running'}));
    //     }, 1500);

    // }),

    // restartProcess: (value) => set(async (state) => {

    //     //set(() => ({processStatus: 'stopping...'}));
    //     let uri = `${rsUrl}/processes/stop/${value}`;
    //     let result = await fetch(uri, {mode: 'cors', method: 'GET'});

    //     //set(() => ({processStatus: 'starting...'}));
    //     uri = `${rsUrl}/processes/start/${value}`;
    //     result = await fetch(uri, {mode: 'cors', method: 'GET'});

    //     setTimeout(async () => {
    //         const pInfo = await fetch(`${rsUrl}/processinfo`, {mode: 'cors'});
    //         const data = await pInfo.json();
    //         set(() => ({processData: data})); //, processStatus:'running'}));
    //     }, 1500);

    // }),


}))

const dataApi = {

    getGameData: async (game) => {

        const uri = `${apiUrl}/games/${game}`;
        console.log('Get Game: ', uri);
        const gameInfo = await fetch(uri, { mode: 'cors', method: 'GET' });
        const data = await gameInfo.json();
        console.log(data);
        return data;
        
    },

    getGames: async () => {

        const games = await fetch(`${apiUrl}/games`, {mode: 'cors'});
        const data = await games.json();

        return data;

    }
}

const unsub = useLnStore.subscribe(state => state.selectedDraw, console.log('HELLO'));


export {
    useLnStore,
    dataApi
}