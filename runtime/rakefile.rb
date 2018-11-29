require_relative '../common.rb'
require_relative './executable/rakefile.rb'
require_relative './io/rakefile.rb'

default runtime: %i[executable io]
