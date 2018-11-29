require_relative '../../common.rb'

path __dir__, 'libio.a' do |t|
  sh 'cargo build --release'
  cp 'target/release/libio.a', t.name
end

default io: File.join(__dir__, 'libio.a')
