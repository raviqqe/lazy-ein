require_relative '../../common.rb'

path __dir__, 'libexecutable.a' do
  sh 'cmake .'
  sh 'make'
end

default executable: File.join(__dir__, 'libexecutable.a')
