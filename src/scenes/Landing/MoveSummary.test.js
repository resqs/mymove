import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { shallow } from 'enzyme';
import MoveSummary from '.';
import store from 'shared/store';

describe('MoveSummary tests', () => {
  let wrapper, div;
  const editMoveFn = jest.fn();
  const resumeMoveFn = jest.fn();
  const entitlementObj = {};
  const serviceMember = {};
  const ordersObj = {};
  const moveObj = { status: 'APPROVED' };
  const ppmObj = { planned_move_date: '2018-06-20' };
  describe('when a move is approved', () => {
    beforeEach(() => {
      wrapper = shallow(
        <MoveSummary
          entitlement={entitlementObj}
          profile={serviceMember}
          orders={ordersObj}
          move={moveObj}
          ppm={ppmObj}
          editMove={editMoveFn}
          resumeMove={resumeMoveFn}
        />,
      );
    });
    it('there is a div', () => {
      // const div = document.createElement('div');
      // ReactDOM.render(
      //   <Provider store={store}>
      //   <MoveSummary
      //     entitlement={entitlementObj}
      //     profile={serviceMember}
      //     orders={ordersObj}
      //     move={moveObj}
      //     ppm={ppmObj}
      //     editMove={editMoveFn}
      //     resumeMove={resumeMoveFn}
      //   />
      // </Provider>,
      //  div,
      // );
      console.log(wrapper.debug());
      console.log(wrapper.contains(<div class="whole_box" />));
      expect(wrapper.find('button').length).toEqual(1);
    });
  });
});
