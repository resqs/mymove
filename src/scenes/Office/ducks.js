import {
  LoadAccountingAPI,
  UpdateAccountingAPI,
  LoadMove,
  LoadOrders,
  LoadServiceMember,
  LoadBackupContacts,
  LoadPPMs,
  ApproveBasics,
} from './api.js';
import * as ReduxHelpers from 'shared/ReduxHelpers';

// Types
const loadAccountingType = 'LOAD_ACCOUNTING';
const updateAccountingType = 'UPDATE_ACCOUNTING';
const loadMoveType = 'LOAD_MOVE';
const loadOrdersType = 'LOAD_ORDERS';
const loadServiceMemberType = 'LOAD_SERVICE_MEMBER';
const loadBackupContactType = 'LOAD_BACKUP_CONTACT';
const loadPPMsType = 'LOAD_PPMS';
const loadDependenciesType = 'LOAD_DEPENDENCIES';
const approveBasicsType = 'APPROVE_BASICS';

const LOAD_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  loadAccountingType,
);

const UPDATE_ACCOUNTING = ReduxHelpers.generateAsyncActionTypes(
  updateAccountingType,
);

const LOAD_MOVE = ReduxHelpers.generateAsyncActionTypes(loadMoveType);

const LOAD_ORDERS = ReduxHelpers.generateAsyncActionTypes(loadOrdersType);

const LOAD_SERVICE_MEMBER = ReduxHelpers.generateAsyncActionTypes(
  loadServiceMemberType,
);

const LOAD_BACKUP_CONTACT = ReduxHelpers.generateAsyncActionTypes(
  loadBackupContactType,
);

const LOAD_PPMS = ReduxHelpers.generateAsyncActionTypes(loadPPMsType);

const LOAD_DEPENDENCIES = ReduxHelpers.generateAsyncActionTypes(
  loadDependenciesType,
);

const APPROVE_BASICS = ReduxHelpers.generateAsyncActionTypes(approveBasicsType);

export const loadAccounting = ReduxHelpers.generateAsyncActionCreator(
  loadAccountingType,
  LoadAccountingAPI,
);

export const updateAccounting = ReduxHelpers.generateAsyncActionCreator(
  updateAccountingType,
  UpdateAccountingAPI,
);

export const loadMove = ReduxHelpers.generateAsyncActionCreator(
  loadMoveType,
  LoadMove,
);

export const loadOrders = ReduxHelpers.generateAsyncActionCreator(
  loadOrdersType,
  LoadOrders,
);

export const loadServiceMember = ReduxHelpers.generateAsyncActionCreator(
  loadServiceMemberType,
  LoadServiceMember,
);

export const loadBackupContacts = ReduxHelpers.generateAsyncActionCreator(
  loadBackupContactType,
  LoadBackupContacts,
);

export const loadPPMs = ReduxHelpers.generateAsyncActionCreator(
  loadPPMsType,
  LoadPPMs,
);

export const approveBasics = ReduxHelpers.generateAsyncActionCreator(
  approveBasicsType,
  ApproveBasics,
);

export function loadMoveDependencies(moveId) {
  const actions = ReduxHelpers.generateAsyncActions(loadDependenciesType);
  return async function(dispatch, getState) {
    dispatch(actions.start());
    try {
      await dispatch(loadMove(moveId));
      const move = getState().office.officeMove;
      await dispatch(loadOrders(move.orders_id));
      const orders = getState().office.officeOrders;
      await dispatch(loadServiceMember(orders.service_member_id));
      const sm = getState().office.officeServiceMember;
      await dispatch(loadBackupContacts(sm.id));
      // TODO: load PPMs in parallel to move using moveId
      await dispatch(loadPPMs(moveId));
      return dispatch(actions.success());
    } catch (ex) {
      return dispatch(actions.error(ex));
    }
  };
}

// Reducer
const initialState = {
  accountingIsLoading: false,
  accountingIsUpdating: false,
  moveIsLoading: false,
  ordersAreLoading: false,
  serviceMemberIsLoading: false,
  backupContactsAreLoading: false,
  ppmsAreLoading: false,
  accountingHasLoadError: null,
  accountingHasLoadSuccess: false,
  accountingHasUpdateError: null,
  accountingHasUpdateSuccess: false,
  moveHasLoadError: null,
  moveHasLoadSuccess: false,
  ordersHaveLoadError: null,
  ordersHaveLoadSuccess: false,
  serviceMemberHasLoadError: null,
  serviceMemberHasLoadSuccess: false,
  backupContactsHaveLoadError: null,
  backupContactsHaveLoadSuccess: false,
  ppmsHaveLoadError: null,
  ppmsHaveLoadSuccess: false,
  loadDependenciesHasError: null,
  loadDependenciesHasSuccess: false,
};

export function officeReducer(state = initialState, action) {
  switch (action.type) {
    case LOAD_ACCOUNTING.start:
      return Object.assign({}, state, {
        accountingIsLoading: true,
        accountingHasLoadSuccess: false,
      });
    case LOAD_ACCOUNTING.success:
      return Object.assign({}, state, {
        accountingIsLoading: false,
        accounting: action.payload,
        accountingHasLoadSuccess: true,
        accountingHasLoadError: false,
      });
    case LOAD_ACCOUNTING.failure:
      return Object.assign({}, state, {
        accountingIsLoading: false,
        accounting: null,
        accountingHasLoadSuccess: false,
        accountingHasLoadError: true,
        error: action.error.message,
      });

    case UPDATE_ACCOUNTING.start:
      return Object.assign({}, state, {
        accountingIsUpdating: true,
        accountingHasUpdateSuccess: false,
      });
    case UPDATE_ACCOUNTING.success:
      return Object.assign({}, state, {
        accountingIsUpdating: false,
        accounting: action.payload,
        accountingHasUpdateSuccess: true,
        accountingHasUpdateError: false,
      });
    case UPDATE_ACCOUNTING.failure:
      return Object.assign({}, state, {
        accountingIsUpdating: false,
        accountingHasUpdateSuccess: false,
        accountingHasUpdateError: true,
        error: action.error.message,
      });

    // Moves
    case LOAD_MOVE.start:
      return Object.assign({}, state, {
        moveIsLoading: true,
        moveHasLoadSuccess: false,
      });
    case LOAD_MOVE.success:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: action.payload,
        moveHasLoadSuccess: true,
        moveHasLoadError: false,
      });
    case LOAD_MOVE.failure:
      return Object.assign({}, state, {
        moveIsLoading: false,
        officeMove: null,
        moveHasLoadSuccess: false,
        moveHasLoadError: true,
        error: action.error.message,
      });

    // ORDERS
    case LOAD_ORDERS.start:
      return Object.assign({}, state, {
        ordersAreLoading: true,
        ordersHaveLoadSuccess: false,
      });
    case LOAD_ORDERS.success:
      return Object.assign({}, state, {
        ordersAreLoading: false,
        officeOrders: action.payload,
        ordersHaveLoadSuccess: true,
        ordersHaveLoadError: false,
      });
    case LOAD_ORDERS.failure:
      return Object.assign({}, state, {
        ordersAreLoading: false,
        officeOrders: null,
        ordersHaveLoadSuccess: false,
        ordersHaveLoadError: true,
        error: action.error.message,
      });

    // SERVICE_MEMBER
    case LOAD_SERVICE_MEMBER.start:
      return Object.assign({}, state, {
        serviceMemberIsLoading: true,
        serviceMemberHasLoadSuccess: false,
      });
    case LOAD_SERVICE_MEMBER.success:
      return Object.assign({}, state, {
        serviceMemberIsLoading: false,
        officeServiceMember: action.payload,
        serviceMemberHasLoadSuccess: true,
        serviceMemberHasLoadError: false,
      });
    case LOAD_SERVICE_MEMBER.failure:
      return Object.assign({}, state, {
        serviceMemberIsLoading: false,
        officeServiceMember: null,
        serviceMemberHasLoadSuccess: false,
        serviceMemberHasLoadError: true,
        error: action.error.message,
      });

    // BACKUP CONTACT
    case LOAD_BACKUP_CONTACT.start:
      return Object.assign({}, state, {
        backupContactsAreLoading: true,
        backupContactsHaveLoadSuccess: false,
      });
    case LOAD_BACKUP_CONTACT.success:
      return Object.assign({}, state, {
        backupContactsAreLoading: false,
        officeBackupContacts: action.payload,
        backupContactsHaveLoadSuccess: true,
        backupContactsHaveLoadError: false,
      });
    case LOAD_BACKUP_CONTACT.failure:
      return Object.assign({}, state, {
        backupContactsAreLoading: false,
        officeBackupContacts: null,
        backupContactsHaveLoadSuccess: false,
        backupContactsHaveLoadError: true,
        error: action.error.message,
      });

    // PPMs
    case LOAD_PPMS.start:
      return Object.assign({}, state, {
        PPMsAreLoading: true,
        PPMsHaveLoadSuccess: false,
      });
    case LOAD_PPMS.success:
      return Object.assign({}, state, {
        PPMsAreLoading: false,
        officePPMs: action.payload,
        PPMsHaveLoadSuccess: true,
        PPMsHaveLoadError: false,
      });
    case LOAD_PPMS.failure:
      return Object.assign({}, state, {
        PPMsAreLoading: false,
        officePPMs: null,
        PPMsHaveLoadSuccess: false,
        PPMsHaveLoadError: true,
        error: action.error.message,
      });

    case APPROVE_BASICS.start:
      return Object.assign({}, state, {
        basicsAreUpdating: true,
      });
    case APPROVE_BASICS.success:
      return Object.assign({}, state, {
        basicsAreUpdating: false,
        officeMove: action.payload,
      });
    case APPROVE_BASICS.failure:
      return Object.assign({}, state, {
        basicsAreUpdating: false,
        error: action.error.message,
      });

    // ALL DEPENDENCIES
    case LOAD_DEPENDENCIES.start:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: false,
        loadDependenciesHasError: false,
      });
    case LOAD_DEPENDENCIES.success:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: true,
        loadDependenciesHasError: false,
      });
    case LOAD_DEPENDENCIES.failure:
      return Object.assign({}, state, {
        loadDependenciesHasSuccess: false,
        loadDependenciesHasError: true,
        error: action.error.message,
      });
    default:
      return state;
  }
}