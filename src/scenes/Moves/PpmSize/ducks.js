import { CreatePpm, UpdatePpm, GetPpm } from './api.js';

// Types
export const SET_PENDING_PPM_SIZE = 'SET_PENDING_PPM_SIZE';
export const CREATE_OR_UPDATE_PPM = 'CREATE_OR_UPDATE_PPM';
export const CREATE_OR_UPDATE_PPM_SUCCESS = 'CREATE_OR_UPDATE_PPM_SUCCESS';
export const CREATE_OR_UPDATE_PPM_FAILURE = 'CREATE_OR_UPDATE_PPM_FAILURE';
export const GET_PPM = 'GET_PPM';
export const GET_PPM_SUCCESS = 'GET_PPM_SUCCESS';
export const GET_PPM_FAILURE = 'GET_PPM_FAILURE';

// Creating ppm
export const createOrUpdatePpmRequest = () => ({
  type: CREATE_OR_UPDATE_PPM,
});

export const createOrUpdatePpmSuccess = item => ({
  type: CREATE_OR_UPDATE_PPM_SUCCESS,
  item,
});

export const createOrUpdatePpmFailure = error => ({
  type: CREATE_OR_UPDATE_PPM_FAILURE,
  error,
});

export const getPpmRequest = () => ({
  type: GET_PPM,
});

export const getPpmSuccess = items => ({
  type: GET_PPM_SUCCESS,
  item: items.length > 0 ? items[0] : null,
});

export const getPpmFailure = error => ({
  type: GET_PPM_FAILURE,
  error,
});

// Action creation
export function setPendingPpmSize(value) {
  return { type: SET_PENDING_PPM_SIZE, payload: value };
}
export function createOrUpdatePpm(moveId, size) {
  return function(dispatch, getState) {
    dispatch(createOrUpdatePpmRequest());
    const state = getState();
    const currentPpm = state.ppm.currentPpm;
    if (currentPpm) {
      console.log('update', currentPpm);
      UpdatePpm(moveId, currentPpm.id, { size })
        .then(item => dispatch(createOrUpdatePpmSuccess(item)))
        .catch(error => dispatch(createOrUpdatePpmFailure(error)));
    } else {
      CreatePpm(moveId, { size })
        .then(item => dispatch(createOrUpdatePpmSuccess(item)))
        .catch(error => dispatch(createOrUpdatePpmFailure(error)));
    }
  };
}

export function loadPpm(moveId) {
  return function(dispatch, getState) {
    dispatch(getPpmRequest());
    const state = getState();
    const currentPpm = state.ppm.currentPpm;
    if (!currentPpm) {
      GetPpm(moveId)
        .then(item => dispatch(getPpmSuccess(item)))
        .catch(error => dispatch(getPpmFailure(error)));
    }
  };
}
// Reducer
const initialState = {
  pendingPpmSize: null,
  currentPpm: null,
  hasSubmitError: false,
  hasSubmitSuccess: false,
};
export function ppmReducer(state = initialState, action) {
  switch (action.type) {
    case SET_PENDING_PPM_SIZE:
      return Object.assign({}, state, {
        pendingPpmSize: action.payload,
      });
    case CREATE_OR_UPDATE_PPM_SUCCESS:
      return Object.assign({}, state, {
        pendingPpmSize: null,
        currentPpm: action.item,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    case CREATE_OR_UPDATE_PPM_FAILURE:
      return Object.assign({}, state, {
        currentPpm: null,
        hasSubmitSuccess: false,
        hasSubmitError: true,
      });
    case GET_PPM_SUCCESS:
      return Object.assign({}, state, {
        currentPpm: action.item,
        hasSubmitSuccess: true,
        hasSubmitError: false,
      });
    default:
      return state;
  }
}
