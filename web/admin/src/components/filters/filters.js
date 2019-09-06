import React from 'react'
import { Input, Radio, Select, DatePicker } from 'antd'

import FilterIcon from '../filter-icon/FilterIcon'
import FilterDropdown from '../filter-dropdown/FilterDropdown'

/* eslint-disable react/display-name */

export const textFilter = (name, { handleFilterSet, handleFilterReset, isFiltered }) => {
  let ref

  return {
    filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => {
      return (
        <FilterDropdown
          onConfirm={handleFilterSet(name, selectedKeys[0], confirm, clearFilters)}
          onReset={handleFilterReset(name, clearFilters)}
        >
          <Input
            ref={el => ref = el}
            value={selectedKeys[0]}
            onChange={e => setSelectedKeys([e.target.value])}
          />
        </FilterDropdown>
      )
    },
    onFilterDropdownVisibleChange: visible => {
      if (visible) {
        setTimeout(() => {
          ref.focus()
        })
      }
    },
    filterIcon: <FilterIcon isFiltered={isFiltered(name)} />,
  }
}

export const radioFilter = (name, values, { handleFilterSet, handleFilterReset, isFiltered }) => {
  return {
    filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => {
      return (
        <FilterDropdown
          onConfirm={handleFilterSet(name, selectedKeys[0], confirm, clearFilters)}
          onReset={handleFilterReset(name, clearFilters)}
        >
          <Radio.Group
            onChange={e => setSelectedKeys([e.target.value])}
            value={selectedKeys[0]}
          >
            {
              values.map(({ value, data }, i) => (<Radio key={i} value={value}>{data}</Radio>))
            }
          </Radio.Group>
        </FilterDropdown>
      )
    },
    filterIcon: <FilterIcon isFiltered={isFiltered(name)} />,
  }
}

export const selectFilter = (name, values, { handleFilterSet, handleFilterReset, isFiltered }) => {
  let ref

  return {
    filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => {
      return (
        <FilterDropdown
          onConfirm={handleFilterSet(name, selectedKeys[0], confirm, clearFilters)}
          onReset={handleFilterReset(name, clearFilters)}
        >
          <Select
            ref={el => ref = el}
            mode='multiple'
            optionFilterProp='children'
            style={{ width: '200px' }}
            onChange={values => setSelectedKeys([values])}
            value={selectedKeys[0]}
            getPopupContainer={triggerNode => triggerNode.parentNode}
          >
            {
              values.map(({ value, data }) => {
                return (
                  <Select.Option key={value}>{data}</Select.Option>
                )
              })
            }
          </Select>
        </FilterDropdown>
      )
    },
    onFilterDropdownVisibleChange: visible => {
      if (visible) {
        setTimeout(() => {
          ref.focus()
        })
      }
    },
    filterIcon: <FilterIcon isFiltered={isFiltered(name)} />,
  }
}

export const dateRangeFilter = (name, { handleFilterSet, handleFilterReset, isFiltered }) => {
  return {
    filterDropdown: ({ setSelectedKeys, selectedKeys, confirm, clearFilters }) => {
      return (
        <FilterDropdown
          onConfirm={handleFilterSet(name, selectedKeys[0], confirm, clearFilters)}
          onReset={handleFilterReset(name, clearFilters)}
        >
          <DatePicker.RangePicker
            onChange={([ from, to ]) => {
              return setSelectedKeys([
                {
                  from: from.toISOString(),
                  to: to.toISOString(),
                },
              ])
            }}
            getCalendarContainer={triggerNode => triggerNode.parentNode}
          />
        </FilterDropdown>
      )
    },
    filterIcon: <FilterIcon isFiltered={isFiltered(name)} />,
  }
}

/* eslint-enable react/display-name */
