import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { shallow } from 'enzyme';
import MoveSummary from '.';

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
      console.log(wrapper.contains(<div class="whole_box" />));
      console.log(
        wrapper
          .find('h2')
          .parents()
          .at(0),
      );
      expect(wrapper.find('button').length).toEqual(1);
    });
  });
});
