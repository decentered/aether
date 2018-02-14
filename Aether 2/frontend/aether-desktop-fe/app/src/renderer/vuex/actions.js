import * as types from './mutation-types'

// export const decrementMain = ({ commit }) => {
//   commit(types.DECREMENT_MAIN_COUNTER)
// }

// export const incrementMain = ({ commit }) => {
//   commit(types.INCREMENT_MAIN_COUNTER)
// }

export const browserBack = ({ commit }) => {
  commit(types.BROWSER_BACK)
}

export const browserForward = ({ commit }) => {
  commit(types.BROWSER_FORWARD)
}
