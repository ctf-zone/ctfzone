import { ApiError } from '../utils/api'

const createReducer = (initial) => (state, { name, action, payload = {} }) => {
  const mixin = { ...initial, ...payload }

  return {
    ...state,
    global: { ...state.global, ...mixin },
    models: {
      ...state.models,
      [name]: { ...state.models[name], ...mixin },
    },
    effects: {
      ...state.effects,
      [name]: {
        ...state.effects[name],
        [action]: { ...state.effects[name][action], ...mixin },
      },
    },
  }
}

const initialState = {
  success: false,
  error: null,
  loading: false,
  time: new Date(),
}

const api = {
  state: {
    global: { ...initialState },
    models: {},
    effects: {},
  },
  reducers: {
    success: createReducer({ success: true }),
    failure: createReducer({ success: false }),
    start: createReducer({ loading: true }),
    stop: createReducer({ loading: false }),
  },
}

export default ({ blacklist = [ 'api' ], whitelist = [] }) => ({
  config: {
    models: {
      api,
    },
  },
  onModel({ name }) {

    if (whitelist.length) {
      if (!whitelist.includes(name)) {
        return
      }
    } else {
      if (blacklist.includes(name)) {
        return
      }
    }

    // Set initial state
    api.state.models[name] = { ...initialState }
    api.state.effects[name] = {}

    Object.keys(this.dispatch[name]).forEach((action) => {

      if (!this.dispatch[name][action].isEffect) {
        return
      }

      api.state.effects[name][action] = { ...initialState }

      const originalEffect = this.dispatch[name][action]

      let wrapperEffect = async(...args) => {
        let effectResult

        try {
          this.dispatch.api.start({ name, action })

          effectResult = await originalEffect(...args)

          this.dispatch.api.success({
            name,
            action,
            payload: { error: null },
          })
        } catch (err) {
          this.dispatch.api.failure({
            name,
            action,
            payload: { error: err },
          })

          if (args.length > 1 && ('throw' in args[args.length - 1])) {
            const { status, data } = err.response
            const { error, errors } = data
            throw new ApiError(status, error, errors)
          }
        } finally {
          this.dispatch.api.stop({
            name,
            action,
            payload: { time: new Date() },
          })
        }

        return effectResult
      }

      wrapperEffect.isEffect = true

      this.dispatch[name][action] = wrapperEffect
    })
  },
})
